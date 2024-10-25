import { useState, useEffect, createContext, ReactNode, useContext, useReducer } from "react";
import Connection from "../libs/socket";
import { OrderListEvent } from "../libs/socket";

interface IWebsocketContext {
    isConnected: boolean;
    orderList: OrderListEvent;
}

export enum ActionType {
    UpdateOrderList,
}

interface UpdateOrderListAction {
    actionType: ActionType.UpdateOrderList;
    payload: OrderListEvent;
}

type Actions = UpdateOrderListAction;

interface State {
    orderList: OrderListEvent;
}

function reducer(state: State, action: Actions) {
    switch (action.actionType) {

        case ActionType.UpdateOrderList:
            return { orderList: action.payload };

        default:
            return state
    }
}

const WebSocketContextProvider = ({ children }: { children: ReactNode | ReactNode[] }) => {
    const [connection, setConnection] = useState({} as Connection);
    const [trigger, setTrigger] = useState(0);
    const [isConnected, setIsConnected] = useState(false);
    const [state, dispatch] = useReducer(reducer, { orderList: { type: "order_list", buy_orders: [], sell_orders: [], buy_quantity: "", sell_quantity: "" } as OrderListEvent });

    useEffect(() => {
        let connection = new Connection();

        connection.onOpen = () => {
            setIsConnected(true);
        };

        connection.onClose = () => {
            setIsConnected(false);
            setTimeout(() => {
                setTrigger(trigger + 1);
            }, 3000);
        };

        connection.onError = () => {
            connection.close();
            setIsConnected(false);
        };

        setConnection(connection);

        return () => { connection.close() }
    }, [trigger]);

    useEffect(() => {
        connection.onOrderList = e => {
            dispatch({ actionType: ActionType.UpdateOrderList, payload: e });
        };
    }, [connection]);

    return (
        <WebsocketContext.Provider value={{ isConnected, orderList: state.orderList }}>
            {children}
        </WebsocketContext.Provider>
    );
};

export const WebsocketContext = createContext<IWebsocketContext>({} as IWebsocketContext);

export function useWebsocket() {
    return useContext(WebsocketContext);
}

export default WebSocketContextProvider;