import { Backdrop, CircularProgress } from "@mui/material";
import { useMe } from "api/me/queries"
import React from "react";

interface IProps {
    children: React.ReactNode;
}
export const GlobalLoading: React.FC<IProps> = ({ children }) => {
    const { isLoading } = useMe();
    return isLoading ? <Backdrop open><CircularProgress /></Backdrop> : children;
}
