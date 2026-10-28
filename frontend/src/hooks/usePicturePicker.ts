import React, { useEffect, useRef, useState } from "react";
import { validatePicture } from "helpers/validatePicture";

// Holds a picture chosen in a file input until it's sent: it validates it
// (helpers/validatePicture), keeps a local preview URL and releases it. Give
// inputRef and pick to the <input type="file">; clear empties everything.
export const usePicturePicker = () => {
  const inputRef = useRef<HTMLInputElement>(null);
  const [picture, setPicture] = useState<File | null>(null);
  const [preview, setPreview] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);

  // Release the local preview URL when it's replaced or the page closes.
  useEffect(() => () => { if (preview) URL.revokeObjectURL(preview); }, [preview]);

  const clear = () => {
    setPicture(null);
    setPreview(null);
    if (inputRef.current) inputRef.current.value = "";
  };

  const pick = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0] ?? null;
    setError(null);
    clear();
    if (!file) return;

    const problem = await validatePicture(file);
    if (problem) {
      setError(problem);
      return;
    }
    setPicture(file);
    setPreview(URL.createObjectURL(file));
  };

  return { inputRef, picture, preview, error, pick, clear };
};
