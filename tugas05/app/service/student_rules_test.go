package service

import (
    "testing"
    "tugas05/app/model"
)

func TestCountTotalPages(t *testing.T) {
    cases := []struct {
        total, limit, want int
    }{
        {0, 10, 0},
        {1, 10, 1},
        {10, 10, 1},
        {11, 10, 2},
        {137, 20, 7},
    }
    for _, tc := range cases {
        got := CountTotalPages(tc.total, tc.limit)
        if got != tc.want {
            t.Errorf("total=%d limit=%d: harap %d, dapat %d", tc.total, tc.limit, tc.want, got)
        }
    }
}

func TestApplyPatch(t *testing.T) {
    initial := model.Student{
        ID:       1,
        NIM:      "S001",
        Name:     "Thaariq",
        Grade:    85.5,
        IsActive: true,
    }

    newActive := false
    req := model.PatchStudentRequest{IsActive: &newActive}
    result, errs := ApplyPatch(initial, req)
    if len(errs) != 0 {
        t.Fatalf("tidak seharusnya ada error: %v", errs)
    }
    if result.IsActive != false {
        t.Error("is_active seharusnya berubah menjadi false")
    }
    if result.Name != "Thaariq" {
        t.Error("field yang tidak dikirim seharusnya tidak berubah")
    }
}

func TestValidateCreate(t *testing.T) {
    req := model.CreateStudentRequest{
        NIM:   "S001",
        Name:  "Valen",
        Grade: 105,
    }
    errs := ValidateCreate(req)
    if len(errs) != 1 {
        t.Errorf("seharusnya 1 error, dapat %d", len(errs))
    }
    if _, ok := errs["grade"]; !ok {
        t.Error("seharusnya error pada grade")
    }
}