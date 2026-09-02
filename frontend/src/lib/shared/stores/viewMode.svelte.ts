export type ViewColumns = 1 | 2 | 3;

let currentColumns = $state<ViewColumns>(
  (() => {
    if (typeof localStorage !== "undefined") {
      const saved = localStorage.getItem("docksea_view_columns");
      if (saved === "1" || saved === "2" || saved === "3") {
        return parseInt(saved, 10) as ViewColumns;
      }
    }
    return 3;
  })(),
);

export const viewModeStore = {
  get columns(): ViewColumns {
    return currentColumns;
  },
  // Retrocompatibilidade
  get mode() {
    return currentColumns === 1 ? "list" : "grid";
  },
  setColumns(cols: ViewColumns) {
    currentColumns = cols;
    if (typeof localStorage !== "undefined") {
      localStorage.setItem("docksea_view_columns", String(cols));
    }
  },
  getGridClass() {
    if (currentColumns === 1) {
      return "grid grid-cols-1 gap-3 w-full";
    }
    if (currentColumns === 2) {
      return "grid grid-cols-2 gap-3 w-full";
    }
    // Force exactly 3 columns regardless of window size
    return "grid grid-cols-3 gap-3 w-full";
  },
};
