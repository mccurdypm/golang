package grafana

import "time"

type GetDetails struct {
    Dashboard struct {
        Annotations struct {
            List []struct {
                BuiltIn    int    `json:"builtIn"`
                Datasource string `json:"datasource"`
                Enable     bool   `json:"enable"`
                Hide       bool   `json:"hide"`
                IconColor  string `json:"iconColor"`
                Name       string `json:"name"`
                Type       string `json:"type"`
            } `json:"list"`
        } `json:"annotations"`
        Editable     bool          `json:"editable"`
        GnetID       interface{}   `json:"gnetId"`
        GraphTooltip int           `json:"graphTooltip"`
        ID           int           `json:"id"`
        Links        []interface{} `json:"links"`
        Panels       []struct {
            Content string `json:"content,omitempty"`
            GridPos struct {
                H int `json:"h"`
                W int `json:"w"`
                X int `json:"x"`
                Y int `json:"y"`
            } `json:"gridPos"`
            ID         int           `json:"id"`
            Links      []interface{} `json:"links"`
            Mode       string        `json:"mode,omitempty"`
            Title      string        `json:"title"`
            Type       string        `json:"type"`
            Columns    []interface{} `json:"columns,omitempty"`
            Datasource string        `json:"datasource,omitempty"`
            FontSize   string        `json:"fontSize,omitempty"`
            PageSize   interface{}   `json:"pageSize,omitempty"`
            Scroll     bool          `json:"scroll,omitempty"`
            ShowHeader bool          `json:"showHeader,omitempty"`
            Sort       struct {
                Col  int  `json:"col"`
                Desc bool `json:"desc"`
            } `json:"sort,omitempty"`
            Styles []struct {
                Alias       string      `json:"alias"`
                DateFormat  string      `json:"dateFormat,omitempty"`
                Pattern     string      `json:"pattern"`
                Type        string      `json:"type"`
                ColorMode   interface{} `json:"colorMode,omitempty"`
                Colors      []string    `json:"colors,omitempty"`
                Decimals    int         `json:"decimals,omitempty"`
                MappingType int         `json:"mappingType,omitempty"`
                Thresholds  []string    `json:"thresholds,omitempty"`
                Unit        string      `json:"unit,omitempty"`
                ValueMaps   []struct {
                    Text  string `json:"text"`
                    Value string `json:"value"`
                } `json:"valueMaps,omitempty"`
            } `json:"styles,omitempty"`
            Targets []struct {
                Dimensions struct {
                    AutoScalingGroupName    string  `json:"AutoScalingGroupName,omitempty"`
                } `json:"dimensions,omitempty"`
                HighResolution bool     `json:"highResolution"`
                MetricName     string   `json:"metricName"`
                Namespace      string   `json:"namespace"`
                Period         string   `json:"period"`
                RefID          string   `json:"refId"`
                Region         string   `json:"region"`
                Statistics     []string `json:"statistics"`
                Target         string   `json:"target"`
                Type           string   `json:"type"`
            } `json:"targets,omitempty"`
            Transform       string      `json:"transform,omitempty"`
            CacheTimeout    interface{} `json:"cacheTimeout,omitempty"`
            ColorBackground bool        `json:"colorBackground,omitempty"`
            ColorValue      bool        `json:"colorValue,omitempty"`
            Colors          []string    `json:"colors,omitempty"`
            Format          string      `json:"format,omitempty"`
            Gauge           struct {
                MaxValue         int  `json:"maxValue"`
                MinValue         int  `json:"minValue"`
                Show             bool `json:"show"`
                ThresholdLabels  bool `json:"thresholdLabels"`
                ThresholdMarkers bool `json:"thresholdMarkers"`
            } `json:"gauge,omitempty"`
            Interval     interface{} `json:"interval,omitempty"`
            MappingType  int         `json:"mappingType,omitempty"`
            MappingTypes []struct {
                Name  string `json:"name"`
                Value int    `json:"value"`
            } `json:"mappingTypes,omitempty"`
            MaxDataPoints   int         `json:"maxDataPoints,omitempty"`
            NullPointMode   string      `json:"nullPointMode,omitempty"`
            NullText        interface{} `json:"nullText,omitempty"`
            Postfix         string      `json:"postfix,omitempty"`
            PostfixFontSize string      `json:"postfixFontSize,omitempty"`
            Prefix          string      `json:"prefix,omitempty"`
            PrefixFontSize  string      `json:"prefixFontSize,omitempty"`
            RangeMaps       []struct {
                From string `json:"from"`
                Text string `json:"text"`
                To   string `json:"to"`
            } `json:"rangeMaps,omitempty"`
            Sparkline struct {
                FillColor string `json:"fillColor"`
                Full      bool   `json:"full"`
                LineColor string `json:"lineColor"`
                Show      bool   `json:"show"`
            } `json:"sparkline,omitempty"`
            TableColumn   string `json:"tableColumn,omitempty"`
            Thresholds    string `json:"thresholds,omitempty"`
            ValueFontSize string `json:"valueFontSize,omitempty"`
            ValueMaps     []struct {
                Op    string `json:"op"`
                Text  string `json:"text"`
                Value string `json:"value"`
            } `json:"valueMaps,omitempty"`
            ValueName   string `json:"valueName,omitempty"`
            AliasColors struct {
            } `json:"aliasColors,omitempty"`
            Bars       bool `json:"bars,omitempty"`
            DashLength int  `json:"dashLength,omitempty"`
            Dashes     bool `json:"dashes,omitempty"`
            Fill       int  `json:"fill,omitempty"`
            Legend     struct {
                Avg     bool `json:"avg"`
                Current bool `json:"current"`
                Max     bool `json:"max"`
                Min     bool `json:"min"`
                Show    bool `json:"show"`
                Total   bool `json:"total"`
                Values  bool `json:"values"`
            } `json:"legend,omitempty"`
            Lines           bool          `json:"lines,omitempty"`
            Linewidth       int           `json:"linewidth,omitempty"`
            Percentage      bool          `json:"percentage,omitempty"`
            Pointradius     int           `json:"pointradius,omitempty"`
            Points          bool          `json:"points,omitempty"`
            Renderer        string        `json:"renderer,omitempty"`
            SeriesOverrides []interface{} `json:"seriesOverrides,omitempty"`
            SpaceLength     int           `json:"spaceLength,omitempty"`
            Stack           bool          `json:"stack,omitempty"`
            SteppedLine     bool          `json:"steppedLine,omitempty"`
            TimeFrom        interface{}   `json:"timeFrom,omitempty"`
            TimeShift       interface{}   `json:"timeShift,omitempty"`
            Tooltip         struct {
                Shared    bool   `json:"shared"`
                Sort      int    `json:"sort"`
                ValueType string `json:"value_type"`
            } `json:"tooltip,omitempty"`
            Xaxis struct {
                Buckets interface{}   `json:"buckets"`
                Mode    string        `json:"mode"`
                Name    interface{}   `json:"name"`
                Show    bool          `json:"show"`
                Values  []interface{} `json:"values"`
            } `json:"xaxis,omitempty"`
            Yaxes []struct {
                Format  string      `json:"format"`
                Label   interface{} `json:"label"`
                LogBase int         `json:"logBase"`
                Max     interface{} `json:"max"`
                Min     interface{} `json:"min"`
                Show    bool        `json:"show"`
            } `json:"yaxes,omitempty"`
            Yaxis struct {
                Align      bool        `json:"align"`
                AlignLevel interface{} `json:"alignLevel"`
            } `json:"yaxis,omitempty"`
            Transparent bool `json:"transparent,omitempty"`
        } `json:"panels"`
        Refresh       string        `json:"refresh"`
        SchemaVersion int           `json:"schemaVersion"`
        Style         string        `json:"style"`
        Tags          []interface{} `json:"tags"`
        Templating    struct {
            List []interface{} `json:"list"`
        } `json:"templating"`
        Time struct {
            From string `json:"from"`
            To   string `json:"to"`
        } `json:"time"`
        Timepicker struct {
            RefreshIntervals []string `json:"refresh_intervals"`
            TimeOptions      []string `json:"time_options"`
        } `json:"timepicker"`
        Timezone string `json:"timezone"`
        Title    string `json:"title"`
        UID      string `json:"uid"`
        Version  int    `json:"version"`
    } `json:"dashboard"`
    Meta struct {
        CanAdmin    bool      `json:"canAdmin"`
        CanEdit     bool      `json:"canEdit"`
        CanSave     bool      `json:"canSave"`
        CanStar     bool      `json:"canStar"`
        Created     time.Time `json:"created"`
        CreatedBy   string    `json:"createdBy"`
        Expires     time.Time `json:"expires"`
        FolderID    int       `json:"folderId"`
        FolderTitle string    `json:"folderTitle"`
        FolderURL   string    `json:"folderUrl"`
        HasAcl      bool      `json:"hasAcl"`
        IsFolder    bool      `json:"isFolder"`
        Provisioned bool      `json:"provisioned"`
        Slug        string    `json:"slug"`
        Type        string    `json:"type"`
        Updated     time.Time `json:"updated"`
        UpdatedBy   string    `json:"updatedBy"`
        URL         string    `json:"url"`
        Version     int       `json:"version"`
    } `json:"meta"`
}
