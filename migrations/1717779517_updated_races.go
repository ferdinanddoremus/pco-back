package migrations

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models/schema"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("24xq4zj5fruqgx6")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("nblf5bll")

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("24xq4zj5fruqgx6")
		if err != nil {
			return err
		}

		// add
		del_categories := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "nblf5bll",
			"name": "categories",
			"type": "select",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"maxSelect": 6,
				"values": [
					"access1",
					"access2",
					"access3",
					"open1",
					"open2",
					"open3"
				]
			}
		}`), del_categories); err != nil {
			return err
		}
		collection.Schema.AddField(del_categories)

		return dao.SaveCollection(collection)
	})
}
