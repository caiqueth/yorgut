// A friend is rated 0 to RATING_MAX on each trait (the API's limit too).
export const RATING_MAX = 3;

// The profile header shows each average as up to this many icons.
export const RATING_ICONS_SHOWN = 10;

// How many icons stand for an average rating (0 to RATING_MAX) in the profile
// header: a full average shows all RATING_ICONS_SHOWN, rounded to whole icons.
export const ratingIcons = (average: number) => Math.round((average / RATING_MAX) * RATING_ICONS_SHOWN);

// An average rating as a share of the maximum: "90%".
export const ratingPercent = (average: number) => `${Math.round((average / RATING_MAX) * 100)}%`;
