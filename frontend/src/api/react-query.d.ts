import "@tanstack/react-query";

declare module "@tanstack/react-query" {
  interface Register {
    mutationMeta: {
      // The form that runs this mutation shows its error itself, so App.tsx
      // doesn't show it again at the bottom of the screen.
      showsErrorInline?: boolean;
    };
  }
}
