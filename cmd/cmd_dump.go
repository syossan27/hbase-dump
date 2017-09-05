package cmd

import (
	"github.com/urfave/cli"
	"hbase-dump/common"
	"github.com/tsuna/gohbase"
	"github.com/tsuna/gohbase/hrpc"
	"context"
	"encoding/json"
	"os"
	"bufio"
	"fmt"
	"github.com/k0kubun/pp"
)

func CmdDump(c *cli.Context) error {
	// Validation
	accessPoint := c.String("access")
	if accessPoint == "" {
		common.Fatal("Error: option 'access' requires an argument")
	}

	table := c.String("table")
	if table == "" {
		common.Fatal("Error: option 'table' requires an argument")
	}

	index := c.String("index")
	if index == "" {
		common.Fatal("Error: option 'index' requires an argument")
	}

	child := c.Bool("child")

	// Get hbase row data
	client := gohbase.NewClient(accessPoint)
	scanRequest, err := hrpc.NewScanStr(context.Background(), table)
	if err != nil {
		common.Fatal("Error: Failed create scan request")
	}

	scanResult := client.Scan(scanRequest)
	var rows *[]map[string]string = new([]map[string]string)
	for {
		r, err := scanResult.Next()
		if err != nil {
			break
		}

		row := map[string]string{}
		for _, cell := range r.Cells {
			row["rowKey"] = string(cell.Row)
			row[string(cell.Qualifier)] = string(cell.Value)
		}

		*rows = append(*rows, row)
	}

	f, err := os.OpenFile(table + ".json", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println(err)
	}
	var writer *bufio.Writer
	writer = bufio.NewWriter(f)

	for _, j := range *rows {
		fileInfo, err := f.Stat()
		if err != nil {
			common.Fatal("Error: あかんやん")
		}

		if fileInfo.Size() > 104857600 {
			f, err := os.OpenFile(table + ".json", os.O_WRONLY|os.O_CREATE, 0600)
			if err != nil {
				fmt.Println(err)
			}
			writer = bufio.NewWriter(f)
		}

		jsonLine, err := json.Marshal(j)
		if err != nil {
			common.Fatal("Error: あかんやん")
		}

		// テーブルがparentかchildかで対応したレコードのヘッダーを生成
		// ES的に必要なアレ
		var recordHeader string
		if child {
			recordHeader = "{\"index\": {\"_index\": \"" + index + "\", \"_type\": \"" + table + "\", \"parent\": \"" + j["id"] + "\"}}\n"
		} else {
			recordHeader = "{\"index\": {\"_index\": \"" + index + "\", \"_type\": \"" + table + "\", \"_id\": \"" + j["id"] + "\"}}\n"
		}

		writer.WriteString(recordHeader)
		writer.Write(jsonLine)
		writer.WriteString("\n")
		writer.Flush()
	}

	/*
	f, err := os.OpenFile(table + ".json", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		fmt.Println(err)
	}
	writer := bufio.NewWriter(f)

	for _, j := range *rows {
		jsonLine, err := json.Marshal(j)
		if err != nil {
			common.Fatal("Error: あかんやん")
		}

		// テーブルがparentかchildかで対応したレコードのヘッダーを生成
		// ES的に必要なアレ
		var recordHeader string
		if child {
			recordHeader = "{\"index\": {\"_index\": \"" + index + "\", \"_type\": \"" + table + "\", \"parent\": \"" + j["id"] + "\"}}\n"
		} else {
			recordHeader = "{\"index\": {\"_index\": \"" + index + "\", \"_type\": \"" + table + "\", \"_id\": \"" + j["id"] + "\"}}\n"
		}

		writer.WriteString(recordHeader)
		writer.Write(jsonLine)
		writer.WriteString("\n")
		writer.Flush()
	}
	*/

	common.Success("Success: Dump data.")

	return nil
}
