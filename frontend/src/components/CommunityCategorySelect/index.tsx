import { MenuItem, TextField } from "@mui/material";
import { CommunityCategory } from "api/generated/model";
import { communityCategoryLabels, communityCategoryOptions } from "helpers/communityCategories";
import React from "react";

interface IProps {
  label?: string;
  // "" while nothing is picked.
  value: CommunityCategory | "";
  onChange: (category: CommunityCategory) => void;
}
// Picks one of the fixed community categories (creating and editing a
// community). Required: the form can't be sent without one.
export const CommunityCategorySelect: React.FC<IProps> = ({ label, value, onChange }) => (
  <TextField
    select
    required
    label={label}
    value={value}
    onChange={(e) => onChange(e.target.value as CommunityCategory)}
    sx={{ minWidth: 260 }}
  >
    {communityCategoryOptions.map((category) => (
      <MenuItem key={category} value={category}>
        {communityCategoryLabels[category]}
      </MenuItem>
    ))}
  </TextField>
);
