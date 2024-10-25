import { useCallback, useEffect, useState } from "react";
import { AxiosResponse } from "axios";

function ButtonInput({ actionText, onClick: _onClick }: { actionText: string, onClick: (price: string, quantity: string) => Promise<AxiosResponse<any, any>> }) {
    const [price, setPrice] = useState("0");
    const [quantity, setQuantity] = useState("0");
    const [disable, setDisable] = useState(true);

    useEffect(() => {
        let priceInt = parseInt(price, 10);
        let quantityInt = parseInt(quantity, 10);
        if (/^\d+$/.test(price) && /^\d+$/.test(quantity) && priceInt >= 1 && priceInt <= 99999 && quantityInt >= 1 && quantityInt <= 99999) {
            setDisable(false);
        } else {
            setDisable(true);
        }
    }, [price, quantity]);

    const onChange = useCallback((e: React.ChangeEvent<HTMLInputElement>, set: React.Dispatch<React.SetStateAction<string>>) => {
        if (e.target.value.length >= 6) return;
        set(e.target.value);
    }, []);

    const onClick = useCallback((price: string, quantity: string) => {
        setDisable(true);
        _onClick(price, quantity)
            .then(() => {
                setPrice("");
                setQuantity("");
            })
            .finally(() => setDisable(false));
    }, []);

    return (
        <>
            <button className={`w-[250px] h-[60px] text-3xl bg-gray-700 disabled:bg-gray-800 disabled:cursor-not-allowed ${!disable && "hover:bg-gray-500 active:bg-gray-600"}`} onClick={() => onClick(price, quantity)} disabled={disable}>{actionText}...</button>
            <div>
                <table className="w-[260px] table-fixed border-separate border-spacing-1">
                    <tbody>
                        <tr className="odd:bg-slate-800"><th>PRICE</th><th>QUANTITY</th></tr>
                        <tr className="odd:bg-slate-800 bg-gray-700">
                            <th><input className="w-full bg-gray-700 text-center" type="number" value={price} onChange={e => onChange(e, setPrice)} /></th>
                            <th><input className="w-full bg-gray-700 text-center" type="number" value={quantity} onChange={e => onChange(e, setQuantity)} /></th>
                        </tr>
                    </tbody>
                </table>
            </div>
        </>
    )
}

export default ButtonInput;