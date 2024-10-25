import { Order } from "src/libs/socket";

function BuyRequestText({ orders, quantity }: { orders: Order[], quantity: string }) {
    return (
        <div>
            <span className="mx-1 text-lg">{quantity}</span>
            <span className="mx-1 text-gray-400">requests to buy at</span>
            <span className="mx-1 text-lg">${orders.length > 0 && orders[0].price}</span>
            <span className="mx-1 text-gray-400">or lower</span>
        </div>
    );
}

export default BuyRequestText;