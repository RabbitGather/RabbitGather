//MATCH  (position:Position)-[]-(article:Article)
//  WHERE ds = distance(position.pt, point({longitude: $longitude, latitude: $latitude})) < $radius
//RETURN article, position, ds

MATCH  (position:Position)-[]-(article:Article)
WITH  article ,position ,distance(position.pt, point({longitude:$longitude, latitude:$latitude})) as distance
  WHERE distance< $radius
RETURN position , article , distance;

//
//  MATCH  (position:Position)-[]-(article:Article)
//  WITH  article ,position ,distance(position.pt, point({longitude:145.488875, latitude:54.448655})) as distance
//  WHERE distance< 999999999
//  RETURN position , article , distance;
