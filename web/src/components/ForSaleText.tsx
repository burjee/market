import { Order } from "src/libs/socket";

function ForSaleText({ orders, quantity }: { orders: Order[], quantity: string }) {
    return (
        <div>
            <span className="mx-1 text-lg">{quantity}</span>
            <span className="mx-1 text-gray-400">for sale starting at</span>
            <span className="mx-1 text-lg">${orders.length > 0 && orders[0].price}</span>
        </div>
    );
}

export default ForSaleText;