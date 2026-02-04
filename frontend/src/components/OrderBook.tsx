import React from 'react';
import { OrderBook as OrderBookType } from '../types';

interface OrderBookProps {
  data: OrderBookType;
}

export const OrderBook: React.FC<OrderBookProps> = ({ data }) => {
  // 安全检查：确保 data 及其属性存在
  const bids = data?.bids ?? [];
  const asks = data?.asks ?? [];
  const hasData = bids.length > 0 && asks.length > 0;

  const spread = hasData
    ? (asks[0].price - bids[0].price).toFixed(2)
    : '-';
  const spreadPercent = hasData
    ? (((asks[0].price - bids[0].price) / bids[0].price) * 100).toFixed(2)
    : '0.00';

  return (
    <div className="h-full flex flex-row fin-panel border-l fin-divider overflow-hidden text-xs font-mono select-none">
       {/* Bids Column (Green/Buys) */}
       <div className="flex-1 flex flex-col border-r fin-divider">
          <div className="p-2 border-b fin-divider font-bold text-slate-400 flex justify-between fin-panel-strong">
             <span>买盘 (Bids)</span>
             <span className="text-[10px] font-normal opacity-70">Size</span>
          </div>
          <div className="flex-1 overflow-hidden">
             {bids.slice(0, 15).map((bid, i) => (
                <div key={`bid-${i}`} className="relative flex justify-between px-2 py-0.5 hover:bg-slate-800/50 cursor-crosshair">
                   <div 
                    className="absolute top-0 left-0 bottom-0 bg-green-900/20 transition-all duration-300" 
                    style={{ width: `${Math.min(bid.percent * 5, 100)}%` }}
                  />
                  <span className="text-green-400 relative z-10">{bid.price.toFixed(2)}</span>
                  <span className="text-slate-300 relative z-10">{bid.size}</span>
                </div>
             ))}
          </div>
       </div>

       {/* Middle Spread Info (Vertical Strip) */}
       <div className="w-24 flex flex-col items-center justify-center border-r fin-divider fin-panel-strong z-10 shadow-inner">
           <div className="text-slate-500 text-[10px] uppercase">Spread</div>
           <div className="text-white font-bold my-1">{spread}</div>
           <div className="text-slate-500 text-[10px]">{spreadPercent}%</div>
       </div>

       {/* Asks Column (Red/Sells) */}
       <div className="flex-1 flex flex-col">
          <div className="p-2 border-b fin-divider font-bold text-slate-400 flex justify-between fin-panel-strong">
             <span>卖盘 (Asks)</span>
             <span className="text-[10px] font-normal opacity-70">Size</span>
          </div>
          <div className="flex-1 overflow-hidden">
            {/* Show Asks ascending (lowest ask first) */}
            {asks.slice(0, 15).map((ask, i) => (
                <div key={`ask-${i}`} className="relative flex justify-between px-2 py-0.5 hover:bg-slate-800/50 cursor-crosshair">
                   <div 
                    className="absolute top-0 right-0 bottom-0 bg-red-900/20 transition-all duration-300" 
                    style={{ width: `${Math.min(ask.percent * 5, 100)}%` }} 
                  />
                  <span className="text-red-400 relative z-10">{ask.price.toFixed(2)}</span>
                  <span className="text-slate-300 relative z-10">{ask.size}</span>
                </div>
            ))}
          </div>
       </div>
    </div>
  );
};
