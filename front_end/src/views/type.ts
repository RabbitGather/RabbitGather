export type Position = {
  x: number;
  y: number;
};
export type ArticleChangeEvent = {
  event: string;
  position: Position;
  timestamp: number;
  id: number;
};
export interface Article {
  position: Position;
  timestamp: number;
  id: number;
}
