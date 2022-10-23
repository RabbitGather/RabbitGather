match(article:Article) where article.id =$id
return
  article.title as title,
  article.content as content,
  article.timestamp as timestamp