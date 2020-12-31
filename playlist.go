package v1

type ReqAddPlayList struct {
	UserLink     string //用户link
	Name         string //酷狗 fileName   歌曲名称
	SongHash     string //酷狗 hash
	AlbumAudioId string //酷狗 album_audio_id
	Singer       string //酷狗 singerName  歌手
	TimeLength   int64  // 歌曲长度
	AlbumName    string //专辑名称
}

type ReqDelPlayList struct {
	Guid string
}

type ReqListPlayList struct {
	UserLink string
	Index    int64
	PageSize int64
}

type ResListPlayList struct {
	Rows  []ListPlayList
	Total int64
}

type ListPlayList struct {
	Guid         string
	UserLink     string
	Name         string //酷狗 fileName   歌曲名称
	SongHash     string //酷狗 hash
	AlbumAudioId string //酷狗 album_audio_id
	CreateTime   int64
	Singer       string //酷狗 singerName  歌手
	TimeLength   int64  // 歌曲长度
	AlbumName    string //专辑名称
}
type ReqGetHitSongs struct {
	UserLink string
	Index    int64
	PageSize int64
}

type ResGetHitSongs struct {
	Rows  []HitSongs
	Total int64
}

type HitSongs struct {
	Name         string //酷狗 fileName   歌曲名称
	AlbumAudioId string //酷狗 album_audio_id
	Singer       string //酷狗 singerName  歌手
	SongHash     string //酷狗 hash
	Top          int64  //热门歌曲排名
	Exist        bool   //歌曲是否已经在歌单
	TimeLength   int64  // 歌曲长度
	AlbumName    string //专辑名称
}
