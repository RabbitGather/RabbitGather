update article
    left join article_tag a
    on article.id = a.article_id set title = ?, content = ?
where id = ? and a.tag_id != 1;