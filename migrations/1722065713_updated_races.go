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
		new_area := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "yg6drvpi",
			"name": "area",
			"type": "number",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"min": null,
				"max": null,
				"noDecimal": false
			}
		}`), new_area); err != nil {
			return err
		}
		collection.Schema.AddField(new_area)

		// add
		new_city := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "m3ohaole",
			"name": "city",
			"type": "text",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"min": null,
				"max": null,
				"pattern": ""
			}
		}`), new_city); err != nil {
			return err
		}
		collection.Schema.AddField(new_city)

		// add
		new_category := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "0j65hyiq",
			"name": "category",
			"type": "text",
			"required": false,
			"presentable": false,
			"unique": false,
			"options": {
				"min": null,
				"max": null,
				"pattern": ""
			}
		}`), new_category); err != nil {
			return err
		}
		collection.Schema.AddField(new_category)

		// add
		new_slug := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "op81v5ie",
			"name": "slug",
			"type": "text",
			"required": true,
			"presentable": false,
			"unique": false,
			"options": {
				"min": null,
				"max": null,
				"pattern": ""
			}
		}`), new_slug); err != nil {
			return err
		}
		collection.Schema.AddField(new_slug)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("24xq4zj5fruqgx6")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("yg6drvpi")

		// remove
		collection.Schema.RemoveField("m3ohaole")

		// remove
		collection.Schema.RemoveField("0j65hyiq")

		// remove
		collection.Schema.RemoveField("op81v5ie")

		return dao.SaveCollection(collection)
	})
}
