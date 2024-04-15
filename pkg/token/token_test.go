package token

import (
	"database/sql"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/rpc"
	"strings"
	"testing"
)

func TestName(t *testing.T) {
	client, err := rpc.Dial("tcp", "10.95.2.48:8001")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	var reply string
	err = client.Call("chain.getChainInfo", nil, &reply)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(reply)
}

func Test1(t *testing.T) {
	contractABI := "{\n  \"type\":\"event\",\n  \"name\":\"Transacted\",\n  \"inputs\":[\n    {\"type\":\"address\",\"name\":\"msgSender\",\"internalType\":\"address\",\"indexed\":false},\n    {\"type\":\"address\",\"name\":\"otherSigner\",\"internalType\":\"address\",\"indexed\":false},\n    {\"type\":\"bytes32\",\"name\":\"operation\",\"internalType\":\"bytes32\",\"indexed\":false},\n    {\"type\":\"address\",\"name\":\"toAddress\",\"internalType\":\"address\",\"indexed\":false},\n    {\"type\":\"uint256\",\"name\":\"value\",\"internalType\":\"uint256\",\"indexed\":false},\n    {\"type\":\"bytes\",\"name\":\"data\",\"internalType\":\"bytes\",\"indexed\":false}\n  ],\n  \"anonymous\":false \n}"
	json, err := abi.JSON(strings.NewReader(contractABI))
	if err != nil {
		t.Fatal()
	}

}

func Test12(t *testing.T) {
	//打开数据库连接
	db, err := sql.Open("mysql", "bm-noah-test:321e*|Jn3f759GMZ@tcp(10.95.201.24:3306)/blockchains_exchange")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// 测试连接是否正常
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	// 执行 SQL 命令
	rows, err := db.Query("SELECT txid FROM headlog_tx_arbitrum WHERE id > ?", 0)
	if err != nil {
		panic(err.Error())
	}
	defer rows.Close()
	// 处理查询结果
	for rows.Next() {
		var txid string
		err = rows.Scan(&txid)
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("txid: %s\n", txid)
	}

	fmt.Println("Connected to MySQL database!")
}
