# Color palettes

`charts/colors` ships a catalog of named color palettes — **categorical**,
**sequential**, and **diverging** — that any chart can use via its `Colors`
prop. Categorical palettes cycle discrete colors across series/slices;
sequential and diverging palettes are value gradients (and are sampled into
discrete colors when applied to a categorical chart).

Browse them visually on the demo app's `/palettes` page (`make run-demo`).

## Choosing a palette

```go
import "github.com/geoffjay/templ-charts/charts/colors"

// Ergonomic, typed, autocomplete-friendly:
bar.BarProps{Colors: colors.Scheme(colors.PaletteTableau10)}

// A colorblind-safe palette:
bar.BarProps{Colors: colors.Scheme(colors.PaletteOkabeIto)}

// An explicit custom color list:
bar.BarProps{Colors: colors.PaletteColors("#4269d0", "#efb118", "#ff725c")}
```

Charts default to `nivo` when `Colors` is unset. Every palette below also has a
typed `PaletteID` constant (e.g. `colors.PaletteObservable10`).

## Discovery API

```go
colors.Palettes()                            // full ordered catalog + metadata
colors.PalettesByKind(colors.KindSequential) // filter by kind
colors.ColorblindSafePalettes()              // accessibility-friendly subset
p, ok := colors.LookupPalette(colors.PaletteSunset)
p.Swatch(8)                                  // preview colors (samples gradients)
p.Ordinal() / p.Sequential() / p.Diverging() // chart-ready configs
```

Sequential/diverging rows below show a 7-stop sample of the gradient.

## Categorical

| id | name | source | colorblind-safe | colors |
|---|---|---|:---:|---|
| `nivo` | Nivo | nivo |  | `#e8c1a0` `#f47560` `#f1e15b` `#e8a838` `#61cdbb` `#97e3d5` |
| `category10` | Category 10 | d3 |  | `#1f77b4` `#ff7f0e` `#2ca02c` `#d62728` `#9467bd` `#8c564b` `#e377c2` `#7f7f7f` `#bcbd22` `#17becf` |
| `accent` | Accent | d3 |  | `#7fc97f` `#beaed4` `#fdc086` `#ffff99` `#386cb0` `#f0027f` `#bf5b17` `#666666` |
| `dark2` | Dark 2 | d3 |  | `#1b9e77` `#d95f02` `#7570b3` `#e7298a` `#66a61e` `#e6ab02` `#a6761d` `#666666` |
| `paired` | Paired | d3 |  | `#a6cee3` `#1f78b4` `#b2df8a` `#33a02c` `#fb9a99` `#e31a1c` `#fdbf6f` `#ff7f00` `#cab2d6` `#6a3d9a` `#ffff99` `#b15928` |
| `pastel1` | Pastel 1 | d3 |  | `#fbb4ae` `#b3cde3` `#ccebc5` `#decbe4` `#fed9a6` `#ffffcc` `#e5d8bd` `#fddaec` `#f2f2f2` |
| `pastel2` | Pastel 2 | d3 |  | `#b3e2cd` `#fdcdac` `#cbd5e8` `#f4cae4` `#e6f5c9` `#fff2ae` `#f1e2cc` `#cccccc` |
| `set1` | Set 1 | d3 |  | `#e41a1c` `#377eb8` `#4daf4a` `#984ea3` `#ff7f00` `#ffff33` `#a65628` `#f781bf` `#999999` |
| `set2` | Set 2 | d3 |  | `#66c2a5` `#fc8d62` `#8da0cb` `#e78ac3` `#a6d854` `#ffd92f` `#e5c494` `#b3b3b3` |
| `set3` | Set 3 | d3 |  | `#8dd3c7` `#ffffb3` `#bebada` `#fb8072` `#80b1d3` `#fdb462` `#b3de69` `#fccde5` `#d9d9d9` `#bc80bd` `#ccebc5` `#ffed6f` |
| `tableau10` | Tableau 10 | d3 |  | `#4e79a7` `#f28e2c` `#e15759` `#76b7b2` `#59a14f` `#edc949` `#af7aa1` `#ff9da7` `#9c755f` `#bab0ab` |
| `observable10` | Observable 10 | curated |  | `#4269d0` `#efb118` `#ff725c` `#6cc5b0` `#3ca951` `#ff8ab7` `#a463f2` `#97bbf5` `#9c6b4e` `#9498a0` |
| `tableau20` | Tableau 20 | curated |  | `#4e79a7` `#a0cbe8` `#f28e2b` `#ffbe7d` `#59a14f` `#8cd17d` `#b6992d` `#f1ce63` `#499894` `#86bcb6` `#e15759` `#ff9d9a` `#79706e` `#bab0ab` `#d37295` `#fabfd2` `#b07aa1` `#d4a6c8` `#9d7660` `#d7b5a6` |
| `material` | Material | curated |  | `#f44336` `#e91e63` `#9c27b0` `#673ab7` `#3f51b5` `#2196f3` `#03a9f4` `#00bcd4` `#009688` `#4caf50` `#8bc34a` `#cddc39` `#ffeb3b` `#ffc107` `#ff9800` `#ff5722` `#795548` `#607d8b` |
| `okabe_ito` | Okabe–Ito | curated | ✓ | `#e69f00` `#56b4e9` `#009e73` `#f0e442` `#0072b2` `#d55e00` `#cc79a7` `#999999` |
| `tol_vibrant` | Tol Vibrant | curated | ✓ | `#0077bb` `#33bbee` `#009988` `#ee7733` `#cc3311` `#ee3377` `#bbbbbb` |
| `tol_muted` | Tol Muted | curated | ✓ | `#332288` `#88ccee` `#44aa99` `#117733` `#999933` `#ddcc77` `#cc6677` `#882255` `#aa4499` `#dddddd` |

## Sequential

| id | name | source | colorblind-safe | colors |
|---|---|---|:---:|---|
| `viridis` | Viridis | d3 | ✓ | `#440154` `#453781` `#33638d` `#21918c` `#32b67a` `#84d44b` `#ece51b` |
| `magma` | Magma | d3 | ✓ | `#000004` `#070661` `#260fc8` `#6b1fdd` `#b637b9` `#ee5e7e` `#fea265` |
| `inferno` | Inferno | d3 | ✓ | `#000004` `#070661` `#260fc8` `#6b1fdd` `#b637b9` `#ee5e7e` `#fea265` |
| `plasma` | Plasma | d3 | ✓ | `#0d0887` `#1b068d` `#260592` `#300596` `#3a049a` `#45039e` `#5003a2` |
| `turbo` | Turbo | d3 |  | `#23171b` `#3987f9` `#2ee5ae` `#95fb51` `#feb927` `#e54813` `#900c00` |
| `cividis` | Cividis | d3 | ✓ | `#002051` `#1f3e6e` `#575c6e` `#7f7c75` `#a49d78` `#d5c164` `#fdea45` |
| `warm` | Warm | d3 |  | `#4d1b3b` `#7b1e5f` `#a92183` `#d724a8` `#fc27cd` `#ff2af1` `#ff42ff` |
| `cool` | Cool | d3 |  | `#3b8ca8` `#3b6da8` `#3b4fa8` `#3b30a8` `#3b12a8` `#3b009b` `#3b0076` |
| `sunset` | Sunset | curated |  | `#f3e79b` `#f9c487` `#f6a180` `#e88087` `#cb6893` `#9c5a9f` `#5c53a5` |
| `teal` | Teal | curated |  | `#d1eeea` `#a9dad9` `#86c4c9` `#69abb8` `#5090a5` `#3c738e` `#2a5674` |
| `emerald` | Emerald | curated |  | `#d3f2a3` `#9ade96` `#6ebf8b` `#4a9c82` `#257a77` `#115a65` `#074050` |

## Diverging

| id | name | source | colorblind-safe | colors |
|---|---|---|:---:|---|
| `red_blue` | Red–Blue | d3 |  | `#67001f` `#b13339` `#f2a88c` `#f2efee` `#bcdae9` `#468fc1` `#053061` |
| `red_yellow_blue` | Red–Yellow–Blue | d3 |  | `#a50026` `#d84133` `#fcb169` `#faf8c1` `#cfe9e4` `#70a5cd` `#313695` |
| `brown_blueGreen` | Brown–BlueGreen | d3 |  | `#543005` `#95601b` `#e0c283` `#edf1ea` `#b1ded8` `#33928a` `#003c30` |
| `spectral` | Spectral | d3 |  | `#9e0142` `#d54649` `#fcb168` `#fbf8b0` `#d1eca6` `#64b7ab` `#5e4fa2` |
| `temps` | Temps | curated |  | `#009392` `#40b087` `#98ca8a` `#ddd793` `#ecb47e` `#e58575` `#cf597e` |

## Notes

- **Colorblind-safe** palettes (Okabe–Ito, Tol Vibrant, Tol Muted, and the
  perceptually-uniform sequential schemes Viridis/Magma/Inferno/Plasma/Cividis)
  are designed to remain distinguishable for common forms of color vision
  deficiency. They're the recommended default for accessibility.
- **Known issue:** the vendored `magma`/`inferno` interpolators currently share
  the same color ramp, so they render identically (and neither exactly matches
  d3's ramp). Tracking a fix to transcribe d3-scale-chromatic's ramps. Prefer
  `viridis` or `cividis` for accurate perceptually-uniform sequential coloring
  in the meantime.
