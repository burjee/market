import { WSURL } from "@/env";

export enum EventType {
    OrderList = "order_list",
}

export interface OrderListEvent {
    type: EventType.OrderList;
    buy_orders: Order[];
    sell_orders: Order[];
    buy_quantity: string;
    sell_quantity: string;
}

export interface Order {
    price: string;
    quantity: string;
}

export type WSEvent = OrderListEvent;

class Connection {
    socket!: WebSocket;

    onOpen: (event: globalThis.Event) => void = () => { };
    onClose: (event: CloseEvent) => void = () => { };
    onOrderList: (e: OrderListEvent) => void = () => { };
    onError: (event: globalThis.Event) => void = () => { };

    constructor() {
        this.open();
    }

    open() {
        this.socket = new WebSocket(WSURL);
        this.socket.onerror = (event) => this.onError(event);
        this.socket.onopen = (event) => this.onOpen(event);
        this.socket.onclose = (event) => this.onClose(event);
        this.socket.onmessage = (event: MessageEvent<any>) => {
            if (event.data === "ping") {
                this.send("pong");
                return;
            }

            let e: WSEvent = JSON.parse(event.data);
            switch (e.type) {
                case EventType.OrderList: { this.onOrderList(e); break; }
            }
        };
    }

    send(message: string) {
        this.socket.send(message);
    }

    close() {
        this.socket.close();
    }
}

export default Connection;