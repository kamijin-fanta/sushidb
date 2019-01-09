package fetcher

type Row struct {
	MetricKey []byte
	TimeStamp int64
	Key       []byte
	Value     []byte
}

type FetchItem struct {
	MetricKey          []byte
	Rows               []Row
	ReadPointTimeStamp int64 // Holds the address of the first element
	ReadPointIndex     int   // Index number of the first element
	ReadCount          int   // Total number of reads per MetricKey
	Stop               bool  // if true, there are no elements after
}

type Resource interface {
	fetch(key []byte, timestamp int64, asc bool) (rows []Row, error error)
}

type Fetcher struct {
	Resource       Resource
	TargetKeys     [][]byte
	Items          []FetchItem
	MaybeHasNext   bool
	Asc            bool
	LimitTimeStamp int64
}

func NewFetcher(metricKeys [][]byte, startTimeStamp int64, limitTimestamp int64, asc bool, resource Resource) Fetcher {
	items := make([]FetchItem, len(metricKeys))

	for i, key := range metricKeys {
		items[i] = FetchItem{
			MetricKey:          key,
			ReadPointTimeStamp: startTimeStamp,
		}
	}
	return Fetcher{
		Resource:       resource,
		TargetKeys:     metricKeys,
		Items:          items,
		MaybeHasNext:   true,
		Asc:            asc,
		LimitTimeStamp: limitTimestamp,
	}
}

func (f *Fetcher) Next(limit int) (rows []Row, error error) {
	size := 0
	for limit > size {
		var near *FetchItem

		for idx := range f.Items {
			item := &f.Items[idx]
			// fetch next items
			if item.Stop == false && len(item.Rows) <= item.ReadPointIndex {
				rows, err := f.Resource.fetch(item.MetricKey, item.ReadPointTimeStamp, f.Asc)
				if err != nil {
					return nil, err
				}
				item.Rows = rows
				if len(item.Rows) > 0 {
					item.ReadPointTimeStamp = item.Rows[0].TimeStamp
				}
				item.ReadPointIndex = 0
				if len(rows) == 0 {
					item.Stop = true
				}
			}

			if item.Stop == false && (near == nil ||
				(f.Asc && near.ReadPointTimeStamp > item.ReadPointTimeStamp) ||
				(!f.Asc && near.ReadPointTimeStamp < item.ReadPointTimeStamp)) {
				near = item
			}
		}

		if near != nil {
			latest := near.Rows[near.ReadPointIndex]
			if f.LimitTimeStamp != 0 && (f.Asc && latest.TimeStamp >= f.LimitTimeStamp || !f.Asc && latest.TimeStamp < f.LimitTimeStamp) {
				f.MaybeHasNext = false
				break
			}
			rows = append(rows, latest)
			near.ReadPointIndex++
			near.ReadCount++
			size++
			if len(near.Rows) > near.ReadPointIndex {
				near.ReadPointTimeStamp = near.Rows[near.ReadPointIndex].TimeStamp
			}
		} else {
			f.MaybeHasNext = false
			break
		}
	}
	return rows, nil
}
