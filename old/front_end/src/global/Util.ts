export function sleep(ms: number) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}
const htmlfontsize = parseInt(
  getComputedStyle(document.getElementsByTagName("html")[0], null).fontSize
);

export function remToPx(rem: number): number {
  return rem * htmlfontsize;
}
