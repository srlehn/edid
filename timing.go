package edid

/*
https://git.linuxtv.org/edid-decode.git/tree/calc-gtf-cvt.cpp
*/

// https://git.linuxtv.org/edid-decode.git/tree/edid-decode.h
//
// Video Timings
// If interlaced is true, then the vertical blanking
// for each field is (vfp + vsync + vbp + 0.5), except for
// the VIC 39 timings that doesn't have the 0.5 constant.
//
// The sequence of the various video parameters is as follows:
//
// border - front porch - sync - back porch - border - active video
//
// Note: this is slightly different from EDID 1.4 which calls
// 'active video' as 'addressable video' and the EDID 1.4 term
// 'active video' includes the borders.
//
// But since borders are rarely used, the term 'active video' will
// typically be the same as 'addressable video', and that's how I
// use it.
type Timings struct {
	// Active horizontal and vertical frame height, excluding any
	// borders, if present.
	// Note: for interlaced formats the active field height is vact / 2
	// unsigned hact, vact;
	HAct uint
	VAct uint
	// unsigned hratio, vratio;
	HRatio uint
	VRatio uint
	// unsigned pixclk_khz;
	PixelClockKHz uint
	// 0: no reduced blanking
	// 1: CVT reduced blanking version 1
	// 2: CVT reduced blanking version 2
	// 2 | RB_ALT: CVT reduced blanking version 2 video-optimized (1000/1001 fps)
	// 3: CVT reduced blanking version 3
	// 3 | RB_ALT: v3 with a horizontal blank of 160
	// 4: GTF Secondary Curve
	// unsigned rb;
	RB uint
	// bool interlaced;
	Interlaced bool
	// The horizontal frontporch may be negative in GTF calculations,
	// so use int instead of unsigned for hfp. Example: 292x176@76.
	// int hfp;
	HFP int
	// unsigned hsync;
	HSync uint
	// The backporch may be negative in buggy detailed timings.
	// So use int instead of unsigned for hbp and vbp.
	// int hbp;
	HBP int
	// bool pos_pol_hsync;
	PosPolHSync bool
	// For interlaced formats the vertical front porch of the Even Field
	// is actually a half-line longer.
	// unsigned vfp, vsync;
	VFP   uint
	VSync uint
	// For interlaced formats the vertical back porch of the Odd Field
	// is actually a half-line longer.
	// int vbp;
	VBP int
	// bool pos_pol_vsync;
	PosPolVSync bool
	// unsigned hborder, vborder;
	HBorder uint
	VBorder uint
	// bool even_vtotal; // special for VIC 39
	EvenVTotal bool
	// bool no_pol_vsync; // digital composite signals have no vsync polarity
	NoPolVSync bool
	// unsigned hsize_mm, vsize_mm;
	HSizeMM uint
	VSizeMM uint
	// bool ycbcr420; // YCbCr 4:2:0 encoding
	YCbCr420 bool
}
