export const isStripeRow = (rowNum: number, elementsPerRow: number): boolean =>
  (Math.floor(rowNum / elementsPerRow) + 1) % 2 === 1;
