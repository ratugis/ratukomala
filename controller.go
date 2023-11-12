package pasetobackend

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/aiteung/atdb"
	"github.com/whatsauth/watoken"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GCFPostHandler(PASETOPRIVATEKEYENV, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	var Response Credential
	Response.Status = false
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	var datauser User
	err := json.NewDecoder(r.Body).Decode(&datauser)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
	} else {
		if IsPasswordValid(mconn, collectionname, datauser) {
			Response.Status = true
			tokenstring, err := watoken.Encode(datauser.Username, os.Getenv(PASETOPRIVATEKEYENV))
			if err != nil {
				Response.Message = "Gagal Encode Token : " + err.Error()
			} else {
				Response.Message = "Selamat Datang"
				Response.Token = tokenstring
			}
		} else {
			Response.Message = "Password Salah"
		}
	}

	return GCFReturnStruct(Response)
}

func GCFReturnStruct(DataStuct any) string {
	jsondata, _ := json.Marshal(DataStuct)
	return string(jsondata)
}

func InsertDataGeojson(MongoConn *mongo.Database, colname string, coordinate []float64, name, volume, tipe string) (InsertedID interface{}) {
	req := new(LonLatProperties)
	req.Type = tipe
	req.Coordinates = coordinate
	req.Name = name
	req.Volume = volume

	ins := atdb.InsertOneDoc(MongoConn, colname, req)
	return ins
}

func UpdateDataGeojson(MongoConn *mongo.Database, colname, name, newVolume, newTipe string) error {
	// Filter berdasarkan nama
	filter := bson.M{"name": name}

	// Update data yang akan diubah
	update := bson.M{
		"$set": bson.M{
			"volume": newVolume,
			"tipe":   newTipe,
		},
	}

	// Mencoba untuk mengupdate dokumen
	_, err := MongoConn.Collection(colname).UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	return nil
}

func DeleteDataGeojson(MongoConn *mongo.Database, colname string, name string) (*mongo.DeleteResult, error) {
	filter := bson.M{"name": name}
	del, err := MongoConn.Collection(colname).DeleteOne(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	return del, nil
}

func GCHandlerFunc(Mongostring, dbname, colname string) string {
	koneksyen := SetConnection(Mongostring, dbname)
	datageo := GetAllData(koneksyen, colname)

	jsoncihuy, _ := json.Marshal(datageo)

	return string(jsoncihuy)
}

func GCFPostCoordinate(Mongostring, dbname, colname string, r *http.Request) string {
	req := new(Credents)
	conn := SetConnection(Mongostring, dbname)
	resp := new(LonLatProperties)
	err := json.NewDecoder(r.Body).Decode(&resp)
	if err != nil {
		req.Status = strconv.Itoa(http.StatusNotFound)
		req.Message = "error parsing application/json: " + err.Error()
	} else {
		req.Status = strconv.Itoa(http.StatusOK)
		Ins := InsertDataGeojson(conn, colname,
			resp.Coordinates,
			resp.Name,
			resp.Volume,
			resp.Type)
		req.Message = fmt.Sprintf("%v:%v", "Berhasil Input data", Ins)
	}
	return GCFReturnStruct(req)
}

func GCFUpdateNameGeojson(Mongostring, dbname, colname string, r *http.Request) string {
	req := new(Credents)
	resp := new(LonLatProperties)
	conn := SetConnection(Mongostring, dbname)
	err := json.NewDecoder(r.Body).Decode(&resp)
	if err != nil {
		req.Status = strconv.Itoa(http.StatusNotFound)
		req.Message = "error parsing application/json: " + err.Error()
	} else {
		req.Status = strconv.Itoa(http.StatusOK)
		Ins := UpdateDataGeojson(conn, colname,
			resp.Name,
			resp.Volume,
			resp.Type)
		req.Message = fmt.Sprintf("%v:%v", "Berhasil Update data", Ins)
	}
	return GCFReturnStruct(req)
}

func GCFDeleteDataGeojson(Mongostring, dbname, colname string, r *http.Request) string {
	req := new(Credents)
	resp := new(LonLatProperties)
	conn := SetConnection(Mongostring, dbname)
	err := json.NewDecoder(r.Body).Decode(&resp)
	if err != nil {
		req.Status = strconv.Itoa(http.StatusNotFound)
		req.Message = "error parsing application/json: " + err.Error()
	} else {
		req.Status = strconv.Itoa(http.StatusOK)
		delResult, delErr := DeleteDataGeojson(conn, colname, resp.Name)
		if delErr != nil {
			req.Status = strconv.Itoa(http.StatusInternalServerError)
			req.Message = "error deleting data: " + delErr.Error()
		} else {
			req.Message = fmt.Sprintf("Berhasil menghapus data. Jumlah data terhapus: %v", delResult.DeletedCount)
		}
	}
	return GCFReturnStruct(req)
}

func InsertUser(db *mongo.Database, collection string, userdata User) string {
	hash, _ := HashPassword(userdata.Password)
	userdata.Password = hash
	atdb.InsertOneDoc(db, collection, userdata)
	return "Username : " + userdata.Username + "\nPassword : " + userdata.Password
}
