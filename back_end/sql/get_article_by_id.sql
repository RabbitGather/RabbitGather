select title, content,UNIX_TIMESTAMP(update_time)
from article as a
    left join article_tag t
        on a.id = t.article_id
where id = ? and t.tag_id != 1;