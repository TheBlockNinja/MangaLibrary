## Manga Library

#### Endpoints

<details><summary>/v2/user/api_key (POST)</summary>

```json
{
	"user_id":1
}
```

</details>


<details><summary>/v2/site (POST)</summary>

```json
{
	"site":{
		"name":"SITE NAME",
		"base_url": "https://testurl.com",
		"search_url":"https://testurl.com?q={search_id}",
		"min_age":18
	},
	"web_components":[
		{
			"name":"Pages",
			"tag":"a",
			"attribute":"class",
			"value":"gallerythumb",
			"is_link":true,
			"is_download":false,
			"link_attributes" :"[\"href\"]",
			"parent":0,
			"delay":5,
			"meta_data":"{}",
			"reverse":false
		},
		{
			"name":"Images",
			"tag":"img",
			"attribute":"src",
			"value":"galleries",
			"is_link":true,
			"is_download":true,
			"link_attributes" :"[\"src\",\"src-data\"]",
			"parent":1,
			"delay":2,
			"meta_data":"{}",
			"reverse":false
		},
        {
			"name":"Title",
			"tag":"meta",
			"attribute":"property",
			"value":"title",
			"is_link":false,
			"is_download":false,
			"link_attributes" :"[]",
			"parent":0,
			"delay":0,
			"meta_data":"{\"name\":{\"attribute\":\"content\"}}",
			"reverse":false
		}
    ]
}
```

</details>


<details><summary>/v2/site (GET)</summary>
</details>