select a.title , a.content , ST_AsBinary( b.coords)
from `article` as a
    left join `article_details` as b on a.id = b.article
    left join  `article_tag` as c
        on a.id = c.article_id and c.tag_id = 1
where c.tag_id is null and  a.id = ?;