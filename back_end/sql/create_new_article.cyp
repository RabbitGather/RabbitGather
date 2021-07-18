CREATE
  (:User {name: $username})-
  [:Create {timestamp: timestamp()}]
  ->(:Article {
    id:      $id,
    title:   $title,
    content: $content, timestamp: timestamp()
  })
    -[:CreateAt {timestamp: timestamp()}]
    ->(:Position {time: timestamp(), pt: point({longitude: $longitude, latitude: $latitude})});