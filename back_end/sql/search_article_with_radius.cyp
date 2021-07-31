MATCH  (position:Position)-[]-(article:Article)
WITH  article ,position ,distance(position.pt, point({longitude:$longitude, latitude:$latitude})) as distance
  WHERE distance>$min_radius AND distance< $max_radius
RETURN position , article , distance;
