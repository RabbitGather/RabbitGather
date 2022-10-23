MATCH  (position:Position)-[]-(article:Article)
WITH article, position, distance(position.pt, point({x: $longitude, y: $latitude})) AS distance
  WHERE distance >= $min_radius AND distance <= $max_radius
RETURN position, article, distance;
