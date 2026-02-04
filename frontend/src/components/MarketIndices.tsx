import React from 'react';
import { MarketIndex } from '../types';

interface MarketIndicesProps {
  indices: MarketIndex[];
}

export const MarketIndices: React.FC<MarketIndicesProps> = ({ indices }) => {
  if (!indices || indices.length === 0) {
    return (
      <div className="flex items-center gap-4 text-xs text-slate-500">
        <span>大盘数据加载中...</span>
      </div>
    );
  }

  return (
    <div className="flex items-center gap-3">
      {indices.map((index) => (
        <MarketIndexItem key={index.code} index={index} />
      ))}
    </div>
  );
};

interface MarketIndexItemProps {
  index: MarketIndex;
}

const MarketIndexItem: React.FC<MarketIndexItemProps> = ({ index }) => {
  const isUp = index.change >= 0;
  const colorClass = isUp ? 'text-red-500' : 'text-green-500';
  const sign = isUp ? '+' : '';

  return (
    <div className="flex flex-col items-center px-2 py-1 rounded hover:bg-slate-700/30 transition-colors cursor-default">
      <span className="text-xs text-slate-400 truncate max-w-16">{index.name}</span>
      <span className={`text-sm font-mono font-medium ${colorClass}`}>
        {index.price.toFixed(2)}
      </span>
      <span className={`text-xs font-mono ${colorClass}`}>
        {sign}{index.changePercent.toFixed(2)}%
      </span>
    </div>
  );
};
