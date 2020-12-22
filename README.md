# bookmark-extract

## Extract bookmarks from browser profile

Note: Don't forget the single-quotes around your template string.

### Example Usage

Print urls, one per line:

```
./bin/bookmark-extract -q -input ~/.config/google-chrome/Default/Bookmarks -tpl '{{.URL}}'
```

Print `name,url`

```
./bin/bookmark-extract -q -input ~/.config/google-chrome/Default/Bookmarks -tpl '{{.Name}},{{.URL}},,'
```

Print `<a href="URL">Name</a>`

```
./bin/bookmark-extract -q -input ~/.config/google-chrome/Default/Bookmarks -tpl '<a href="{{.URL}}">{{.Name}}</a>'
```


### available template vars

```
DateAdded    string
GUID         string
ID           string
Name         string
Type         string
URL          string 
DateModified string     
```
