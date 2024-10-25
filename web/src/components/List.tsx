import { Order } from "src/libs/socket";

function List({ orders }: { orders: Order[] }) {
    return (<div className="flex-1 flex justify-center items-start border-t border-slate-600">
        <table className="w-[260px] mt-8 table-fixed border-separate border-spacing-1 text-gray-400">
            <thead>
                <tr><th>Price</th><th>Quantity</th></tr>
            </thead>
            <tbody>
                {
                    orders.map(order => <tr key={order.price} className="odd:bg-slate-800"><th>${order.price}</th><th>{order.quantity}</th></tr>)
                }
            </tbody>
        </table>
    </div>)
}

export default List;