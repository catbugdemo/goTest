package auto_generate

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"testing"
)

func TestDataScrpitGenerate2(t *testing.T) {

}

func TestHToP(t *testing.T) {
	t.Run("initStruct", func(t *testing.T) {
		toP := HToP{}
		toP.initStruct()
		assert.Equal(t, "127.0.0.1",toP.DataSource["host"])
	})
	t.Run("create", func(t *testing.T) {
		toP := HToP{}
		toP.create()

		getwd, _ := os.Getwd()
		exist := isExist(path.Join(getwd,"auto_scrpit", "hdfs_to_pg", "hdfs_to_pg.json"))
		assert.Equal(t, true,exist)
		//os.Remove(path.Join(getwd,"auto_scrpit"))
	})
}

func TestGetDatabase(t *testing.T) {
	gamedb = map[string]interface{}{
		"db":       "postgres",
		"host":     "49.234.137.226",
		"port":     "5432",
		"user":     "zonst_xyx",
		"dbname":   "game_test",
		"sslmode":  "disable",
		"password": "zonst_xyx_fengtao",
	}
	t.Run("getPostgresqlColumn", func(t *testing.T) {
		column, e := getPostgresqlColumn(gamedb, "user_send_process")
		assert.Nil(t, e)
		fmt.Println(column)
	})
	t.Run("getDatabase", func(t *testing.T) {
		database, e := getDatabase(gamedb, "user_send_process")
		assert.Nil(t, e)
		fmt.Println(database)
	})
}

func isExist(path string)(bool){
	_, err := os.Stat(path)
	if err != nil{
		if os.IsExist(err){
			return true
		}
		if os.IsNotExist(err){
			return false
		}
		fmt.Println(err)
		return false
	}
	return true
}