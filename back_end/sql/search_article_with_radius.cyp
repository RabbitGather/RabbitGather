MATCH  (position:Position)-[]-(article:Article)
WITH article, position, distance(position.pt, point({longitude: $longitude, latitude: $latitude})) AS distance
  WHERE distance >= $min_radius AND distance <= $max_radius
RETURN position, article, distance;


//MATCH  (position:Position)-[]-(article:Article)
//WITH article, position,
//     distance(position.pt, point({longitude: 121.51187490970621, latitude: 25.040056717110396})) AS distance
//  WHERE distance < 50
//RETURN position, article, distance;