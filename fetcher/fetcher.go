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
	Resource   Resource
	TargetKeys [][]byte
	Items      []FetchItem
}

func NewFetcher(metricKeys [][]byte, startTimeStamp int64, resource Resource) Fetcher {
	items := make([]FetchItem, len(metricKeys))

	for i, key := range metricKeys {
		items[i] = FetchItem{
			MetricKey:          key,
			ReadPointTimeStamp: startTimeStamp,
		}
	}
	return Fetcher{
		Items:    items,
		Resource: resource,
	}
}

func (f *Fetcher) FetchRecursive(limit int, limitTimestamp int64, asc bool) (rows []Row, error error) {
	size := 0
	for limit > size {
		var near *FetchItem

		for idx := range f.Items {
			item := &f.Items[idx]
			// fetch next items
			if item.Stop == false && len(item.Rows) <= item.ReadPointIndex {
				rows, err := f.Resource.fetch(item.MetricKey, item.ReadPointTimeStamp, asc)
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
				(asc && near.ReadPointTimeStamp > item.ReadPointTimeStamp) ||
				(!asc && near.ReadPointTimeStamp < item.ReadPointTimeStamp)) {
				near = item
			}
		}

		if near != nil {
			latest := near.Rows[near.ReadPointIndex]
			if limitTimestamp != 0 && (asc && latest.TimeStamp >= limitTimestamp || !asc && latest.TimeStamp < limitTimestamp) {
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
			break
		}
	}
	return rows, nil
}
