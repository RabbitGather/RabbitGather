CREATE
  (:User {id: $user_id})-
  [:Create]
  ->(:Article {
    id:      $article_id
  })
    -[:At]
    ->(:Position {pt: point({x: $x, y: $y})});

//CREATE
//  (:User {id:2})-
//  [:Create]
//  ->(:Article {
//    id:     1
//  })
//    -[:At]
//    ->(:Position {pt: point({x:24.211683466965695, y: 120.68043825584145})});


//CREATE
//  (:User {name: '$username'})-
//  [:Create {timestamp: timestamp()}]
//  ->(:Article {
//    id:      9849851,
//    title:  ' $title',
//    content: '$content', timestamp: timestamp()
//  })
//    -[:CreateAt {timestamp: timestamp()}]
//    ->(:Position {time: timestamp(), pt: point({longitude: 121.51187490970621, latitude: 25.040056717110396})});