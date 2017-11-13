select id,userid,title,createdatetime,content
from Article where
{{if gt 1 .count}}
id=?id
{{else}}
userid=?userid
{{end}}
