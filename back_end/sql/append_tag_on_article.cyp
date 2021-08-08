// 在指定文章上加標籤
match(a:Article {id:$article_id})
SET a.tag = case when $the_tag in a.tag then a.tag else a.tag + $the_tag end
return a;






match(a:Article {id:1421778711388098560})
SET a.tag = case when a.tag is null then ['ABC'] when 'ABC' in a.tag then a.tag else a.tag + 'ABC' end
return a;


match(a:Article {id:1421778711388098560})
where 'ABC' in a.tag
return a;



match(a:Article {id:1421778711388098560})
where a is not null
SET a.tag = [t IN a.tag WHERE a <> "ABC"]
return a;