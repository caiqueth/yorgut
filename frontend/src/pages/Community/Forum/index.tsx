import { Box } from "@mui/material";
import React from "react";
import { useParams } from "react-router-dom";

interface IProps { }
export const Forum: React.FC<IProps> = () => {
    const params = useParams();
    const communityId = parseInt(params.communityId ?? "");

    return <Box>Fórum da comunidade {communityId}</Box>
}