import React from 'react';
import {
  ComposedChart,
  LineChart,
  Line,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  Cell
} from 'recharts';
import { KLineData, TimePeriod, Stock } from '../types';

interface StockChartProps {
  data: KLineData[];
  period: TimePeriod;
  onPeriodChange: (p: TimePeriod) => void;
  stock?: Stock;
}

// Custom shape for candlestick
const Candlestick = (props: any) => {
  const { x, y, width, height, low, high, open, close } = props;
  const isGrowing = close > open;
  const color = isGrowing ? '#ef4444' : '#22c55e'; // 红涨绿跌（中国A股标准）

  const unitHeight = height / (high - low);
  
  const bodyTop = y + (high - Math.max(open, close)) * unitHeight;
  const bodyBottom = y + (high - Math.min(open, close)) * unitHeight;
  const bodyLen = Math.max(2, bodyBottom - bodyTop);

  return (
    <g>
      {/* Wick */}
      <line x1={x + width / 2} y1={y} x2={x + width / 2} y2={y + height} stroke={color} strokeWidth={1} />
      {/* Body */}
      <rect
        x={x}
        y={bodyTop}
        width={width}
        height={bodyLen}
        fill={color}
        stroke={color} // Fill the gap
      />
    </g>
  );
};

export const StockChart: React.FC<StockChartProps> = ({ data, period, onPeriodChange, stock }) => {
  // Guard clause for empty data to prevent crashes
  if (!data || data.length === 0) {
    return (
      <div className="h-full w-full fin-panel flex items-center justify-center">
        <span className="text-slate-500 text-sm animate-pulse">加载市场数据中...</span>
      </div>
    );
  }

  // Transform data for the chart
  const chartData = data.map(d => ({
    ...d,
    range: [d.low, d.high] as [number, number],
  }));

  const lastClose = data[data.length - 1].close;
  const isIntraday = period === '1m';

  const periods: { id: TimePeriod; label: string }[] = [
    { id: '1m', label: '分时' },
    { id: '1d', label: '日K' },
    { id: '1w', label: '周K' },
    { id: '1mo', label: '月K' },
  ];

  return (
    <div className="h-full w-full fin-panel flex flex-col">
      {/* Header with Tabs and Info */}
      <div className="flex items-center justify-between px-2 py-1 border-b fin-divider fin-panel-strong z-10">
         <div className="flex gap-1">
            {periods.map((p) => (
              <button
                key={p.id}
                onClick={() => onPeriodChange(p.id)}
                className={`text-xs px-3 py-1 rounded transition-colors ${
                  period === p.id 
                    ? 'bg-slate-800/80 text-cyan-300 font-bold' 
                    : 'text-slate-400 hover:text-slate-200 hover:bg-slate-800/40'
                }`}
              >
                {p.label}
              </button>
            ))}
         </div>
         <div className="text-xs text-slate-400 font-mono flex gap-4">
           {isIntraday ? (
             <>
               <span>现价: <span className="text-cyan-300">{(stock?.price || lastClose).toFixed(2)}</span></span>
               <span>均价: <span className="text-yellow-400">{data[data.length - 1].avg?.toFixed(2) || '--'}</span></span>
               <span>最高: <span className="text-red-400">{(stock?.high || Math.max(...data.map(d => d.high))).toFixed(2)}</span></span>
               <span>最低: <span className="text-green-400">{(stock?.low || Math.min(...data.map(d => d.low))).toFixed(2)}</span></span>
             </>
           ) : (
             <>
               <span>收: <span className="text-cyan-300">{lastClose.toFixed(2)}</span></span>
               <span>开: {data[data.length - 1].open.toFixed(2)}</span>
               <span>高: <span className="text-red-400">{data[data.length - 1].high.toFixed(2)}</span></span>
               <span>低: <span className="text-green-400">{data[data.length - 1].low.toFixed(2)}</span></span>
               {data[data.length - 1].ma5 && (
                 <>
                   <span>MA5: <span className="text-yellow-400">{data[data.length - 1].ma5?.toFixed(2)}</span></span>
                   <span>MA10: <span className="text-purple-400">{data[data.length - 1].ma10?.toFixed(2)}</span></span>
                   <span>MA20: <span className="text-orange-400">{data[data.length - 1].ma20?.toFixed(2)}</span></span>
                 </>
               )}
             </>
           )}
         </div>
      </div>

      <div className="flex-1 min-h-0 relative">
        <ResponsiveContainer width="100%" height="100%">
          {isIntraday ? (
             <LineChart data={chartData} margin={{ top: 10, right: 10, left: 0, bottom: 0 }}>
                <CartesianGrid strokeDasharray="3 3" stroke="#1e293b" vertical={false} />
                <XAxis
                  dataKey="time"
                  tick={{ fill: '#94a3b8', fontSize: 10 }}
                  axisLine={{ stroke: '#334155' }}
                  tickLine={false}
                  minTickGap={30}
                  interval="preserveStartEnd"
                />
                <YAxis
                  domain={['auto', 'auto']}
                  orientation="right"
                  tick={{ fill: '#94a3b8', fontSize: 10 }}
                  axisLine={false}
                  tickLine={false}
                  tickFormatter={(val) => val.toFixed(2)}
                />
                <Tooltip
                  contentStyle={{ backgroundColor: '#0f172a', border: '1px solid #1e293b', color: '#e2e8f0' }}
                  itemStyle={{ color: '#38bdf8' }}
                  labelStyle={{ color: '#94a3b8' }}
                  content={({ active, payload, label }) => {
                    if (!active || !payload || payload.length === 0) return null;
                    const d = payload[0]?.payload;
                    if (!d) return null;
                    return (
                      <div style={{ backgroundColor: '#0f172a', border: '1px solid #1e293b', padding: '8px 12px', borderRadius: '4px' }}>
                        <div style={{ color: '#94a3b8', marginBottom: '4px' }}>{label}</div>
                        <div style={{ color: '#38bdf8' }}>价格: {d.close?.toFixed(2)}</div>
                        {d.avg && <div style={{ color: '#facc15' }}>均价: {d.avg?.toFixed(2)}</div>}
                      </div>
                    );
                  }}
                />
                <Line
                  type="linear"
                  dataKey="close"
                  stroke="#38bdf8"
                  strokeWidth={1.5}
                  dot={false}
                  isAnimationActive={false}
                  name="价格"
                />
                {/* 分时均价线 */}
                <Line
                  type="linear"
                  dataKey="avg"
                  stroke="#facc15"
                  strokeWidth={1}
                  dot={false}
                  isAnimationActive={false}
                  name="均价"
                />
             </LineChart>
          ) : (
            <ComposedChart data={chartData} margin={{ top: 10, right: 10, left: 0, bottom: 0 }}>
              <CartesianGrid strokeDasharray="3 3" stroke="#1e293b" vertical={false} />
              <XAxis 
                dataKey="time" 
                tick={{ fill: '#94a3b8', fontSize: 10 }} 
                axisLine={{ stroke: '#334155' }}
                tickLine={false}
                minTickGap={30}
              />
              <YAxis 
                domain={['auto', 'auto']} 
                orientation="right" 
                tick={{ fill: '#94a3b8', fontSize: 10 }} 
                axisLine={false}
                tickLine={false}
                tickFormatter={(val) => val.toFixed(1)}
              />
              <Tooltip
                contentStyle={{ backgroundColor: '#0f172a', border: '1px solid #1e293b', color: '#e2e8f0' }}
                itemStyle={{ color: '#cbd5f5' }}
                labelStyle={{ color: '#94a3b8' }}
                content={({ active, payload, label }) => {
                  if (!active || !payload || payload.length === 0) return null;
                  const d = payload[0]?.payload;
                  if (!d) return null;
                  return (
                    <div style={{ backgroundColor: '#0f172a', border: '1px solid #1e293b', padding: '8px 12px', borderRadius: '4px' }}>
                      <div style={{ color: '#94a3b8', marginBottom: '4px' }}>{label}</div>
                      <div style={{ color: '#e2e8f0' }}>开: {d.open?.toFixed(2)}</div>
                      <div style={{ color: '#ef4444' }}>高: {d.high?.toFixed(2)}</div>
                      <div style={{ color: '#22c55e' }}>低: {d.low?.toFixed(2)}</div>
                      <div style={{ color: '#38bdf8' }}>收: {d.close?.toFixed(2)}</div>
                      {d.ma5 && <div style={{ color: '#facc15' }}>MA5: {d.ma5?.toFixed(2)}</div>}
                      {d.ma10 && <div style={{ color: '#a855f7' }}>MA10: {d.ma10?.toFixed(2)}</div>}
                      {d.ma20 && <div style={{ color: '#f97316' }}>MA20: {d.ma20?.toFixed(2)}</div>}
                    </div>
                  );
                }}
              />
              <Bar
                dataKey="range"
                shape={<Candlestick />}
                isAnimationActive={false}
              >
                {
                  chartData.map((entry, index) => (
                    <Cell key={`cell-${index}`} fill={entry.close > entry.open ? '#ef4444' : '#22c55e'} />
                  ))
                }
              </Bar>
              {/* 均线 */}
              <Line type="linear" dataKey="ma5" stroke="#facc15" strokeWidth={1} dot={false} isAnimationActive={false} name="MA5" />
              <Line type="linear" dataKey="ma10" stroke="#a855f7" strokeWidth={1} dot={false} isAnimationActive={false} name="MA10" />
              <Line type="linear" dataKey="ma20" stroke="#f97316" strokeWidth={1} dot={false} isAnimationActive={false} name="MA20" />
            </ComposedChart>
          )}
        </ResponsiveContainer>
      </div>
      
      {/* Volume Chart at bottom - Color Coded for Buy/Sell */}
      <div className="h-16 border-t fin-divider">
         <ResponsiveContainer width="100%" height="100%">
          <ComposedChart data={chartData} margin={{ top: 0, right: 10, left: 0, bottom: 0 }}>
             <Bar 
              dataKey="volume" 
              isAnimationActive={false} 
            >
               {chartData.map((entry, index) => (
                  <Cell 
                    key={`vol-${index}`} 
                    fill={entry.close >= entry.open ? '#ef4444' : '#22c55e'} 
                    opacity={0.5} 
                  />
              ))}
            </Bar>
          </ComposedChart>
         </ResponsiveContainer>
      </div>
    </div>
  );
};
