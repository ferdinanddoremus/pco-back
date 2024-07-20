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

		// add
		new_categories := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "czu2eawo",
			"name": "categories",
			"type": "relation",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"collectionId": "uqe9211w9w5w4i8",
				"cascadeDelete": false,
				"minSelect": null,
				"maxSelect": null,
				"displayFields": null
			}
		}`), new_categories); err != nil {
			return err
		}
		collection.Schema.AddField(new_categories)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("24xq4zj5fruqgx6")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("czu2eawo")

		return dao.SaveCollection(collection)
	})
}
