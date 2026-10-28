import { GlobalStateContext } from "contexts/GlobalStateContext";
import { useContext } from "react"

export const useGlobalStateContext = () => {
    const context = useContext(GlobalStateContext);

    if (!context) throw Error("useGlobalStateContext must be used inside an GlobalStateContextProvider");

    return context;
}
