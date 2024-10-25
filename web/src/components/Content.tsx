import { HTTPURL } from "@/env";

import { useCallback } from "react";

import axios from "axios";
import { toast, ToastContainer } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';

import ForSaleText from "./ForSaleText";
import BuyRequestText from "./BuyRequestText";
import ButtonInput from "./ButtonInput";
import List from "./List";
import { useWebsocket } from "../contexts/WebSocketContext";

function newReqBody(type: string, price: string, quantity: string) {
    return {
        type,
        price: parseInt(price, 10),
        quantity: parseInt(quantity, 10)
    };
}

function newReqToast(type: string) {
    return {
        pending: `${type} request is pending`,
        success: `${type} request success`,
        error: 'please try again later',
    }
};

function Content() {
    const { orderList } = useWebsocket();

    const buy = useCallback((price: string, quantity: string) => {
        let url = `${HTTPURL}api/buy`
        let body = newReqBody("buy", price, quantity);
        return toast.promise(axios.post(url, body), newReqToast("buy"));
    }, []);

    const sell = useCallback((price: string, quantity: string) => {
        let url = `${HTTPURL}api/sell`
        let body = newReqBody("sell", price, quantity);
        return toast.promise(axios.post(url, body), newReqToast("sell"));
    }, []);

    return (
        <>
            <div className="w-[800px] h-[600px] bg-slate-900 flex">
                <div className="flex-1 m-5 bg-gray-900 border border-black flex flex-col">
                    <div className="flex-1 flex flex-col justify-evenly items-center">
                        <ForSaleText orders={orderList.sell_orders} quantity={orderList.sell_quantity} />
                        <ButtonInput actionText="Buy" onClick={buy} />
                    </div>
                    <List orders={orderList.sell_orders} />
                </div>
                <div className="flex-1 m-5 bg-gray-900 border border-black flex flex-col">
                    <div className="flex-1 flex flex-col justify-evenly items-center">
                        <BuyRequestText orders={orderList.buy_orders} quantity={orderList.buy_quantity} />
                        <ButtonInput actionText="Sell" onClick={sell} />
                    </div>
                    <List orders={orderList.buy_orders} />
                </div>
            </div>
            <ToastContainer autoClose={3500} limit={8} pauseOnFocusLoss={false} closeOnClick theme="dark" />
        </>
    )
}

export default Content;
