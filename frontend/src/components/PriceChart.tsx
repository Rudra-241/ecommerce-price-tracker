import {
  CartesianGrid,
  Line,
  LineChart,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from 'recharts'
import type { PriceStamp } from '../types'
import { formatDate, formatPrice } from '../lib/format'

export function PriceChart({ history }: { history: PriceStamp[] }) {
  const data = history.map((stamp) => ({
    label: formatDate(stamp.ChangedAt),
    price: stamp.Price,
  }))

  if (data.length === 0) {
    return (
      <div className="grid h-72 place-items-center text-text/50">No price history yet.</div>
    )
  }

  return (
    <div className="h-72 w-full">
      <ResponsiveContainer width="100%" height="100%">
        <LineChart data={data} margin={{ top: 8, right: 12, bottom: 8, left: 8 }}>
          <CartesianGrid stroke="#141c0218" vertical={false} />
          <XAxis
            dataKey="label"
            tick={{ fill: '#141c02', fontSize: 12 }}
            tickLine={false}
            axisLine={{ stroke: '#141c02' }}
            minTickGap={24}
          />
          <YAxis
            tick={{ fill: '#141c02', fontSize: 12 }}
            tickLine={false}
            axisLine={{ stroke: '#141c02' }}
            width={72}
            tickFormatter={(v: number) => formatPrice(v)}
          />
          <Tooltip
            formatter={(v: number) => [formatPrice(v), 'Price']}
            contentStyle={{
              border: '2px solid #141c02',
              borderRadius: 12,
              fontFamily: 'Belanosima',
              background: '#fbfef4',
            }}
          />
          <Line
            type="monotone"
            dataKey="price"
            stroke="#a449f2"
            strokeWidth={3}
            dot={{ r: 3, fill: '#a449f2', stroke: '#141c02', strokeWidth: 1 }}
            activeDot={{ r: 5, fill: '#b9f027', stroke: '#141c02', strokeWidth: 2 }}
          />
        </LineChart>
      </ResponsiveContainer>
    </div>
  )
}
