***Запуск через***
```
# соберет сервер, клиент и запустит сервер и клиент
make
# можно несколько раз запустить созданный клиент
make run-client
```

```go
func performPoW(challenge string, difficulty int) (string, int) {
    nonce := 0
    var hash string
    
    // Используем параметр N, который увеличивается с ростом сложности
    N := 1024 * (1 << uint(difficulty)) // Начальное значение 1024, увеличивается экспоненциально
    r := 8
    p := 1
    
    for {
        nonce++
        record := fmt.Sprintf("%s%d", challenge, nonce)
        h, _ := scrypt.Key([]byte(record), []byte(challenge), N, r, p, 32)
        hash = hex.EncodeToString(h[:])
        
        // Простая проверка: хэш должен начинаться с двух нулей для усложнения задачи
        if strings.HasPrefix(hash, "00") {
            break
        }
    }
    
    return hash, nonce
}
```

Server verify:
```go
func verifyPoW(challenge string, nonce int, hash string, difficulty int) bool {
	record := fmt.Sprintf("%s%d", challenge, nonce)
	N := 1024 * (1 << uint(difficulty)) // Начальное значение 1024, увеличивается экспоненциально
	r := 8
	p := 1
	h, _ := scrypt.Key([]byte(record), []byte(challenge), N, r, p, 32)
	calculatedHash := hex.EncodeToString(h[:])

	return strings.HasPrefix(calculatedHash, "00") && calculatedHash == hash
}
```

Попытка добиться линейного увелчения времени в зависимости от сложности, которая меняется в зависимости от скорости перебора хешей на основе статистики за последние 5 расчетов хэша для защиты от DDOS:
```go
Server is listening on port :12345
2024/07/10 07:34:33 INFO handle conn currentDifficulty=1 nonce=407 size=101
2024/07/10 07:34:35 INFO handle conn currentDifficulty=1 nonce=233 size=101
2024/07/10 07:34:38 INFO handle conn currentDifficulty=1 nonce=280 size=101
2024/07/10 07:34:41 INFO handle conn currentDifficulty=1 nonce=655 size=101
2024/07/10 07:34:45 INFO handle conn currentDifficulty=1 nonce=34 size=100
2024/07/10 07:34:47 INFO handle conn currentDifficulty=2 nonce=186 size=101
2024/07/10 07:34:51 INFO handle conn currentDifficulty=3 nonce=326 size=101
2024/07/10 07:34:56 INFO handle conn currentDifficulty=4 nonce=168 size=101
2024/07/10 07:35:02 INFO handle conn currentDifficulty=5 nonce=83 size=100
2024/07/10 07:35:20 INFO handle conn currentDifficulty=6 nonce=165 size=101
2024/07/10 07:35:31 INFO handle conn currentDifficulty=5 nonce=145 size=101
2024/07/10 07:35:37 INFO handle conn currentDifficulty=4 nonce=104 size=101
2024/07/10 07:35:43 INFO handle conn currentDifficulty=3 nonce=309 size=101
```

при дефолтной цели расчета времени хэша - 500 ms, окна последних расчетов = 5 элементов, сложности - 1
```go
const (
    initialDifficulty   = 1
    targetCalculateTime = 500 * time.Millisecond
    calculateWindow     = 5
)
```

Данные машины на которой выполнялось исследование:
```bash
    $ lscpu
    Architecture:            x86_64
      CPU op-mode(s):        32-bit, 64-bit
      Address sizes:         46 bits physical, 48 bits virtual
      Byte Order:            Little Endian
    CPU(s):                  20
      On-line CPU(s) list:   0-19
    Vendor ID:               GenuineIntel
      Model name:            12th Gen Intel(R) Core(TM) i7-12700K
        CPU family:          6
        Model:               151
        Thread(s) per core:  2
        Core(s) per socket:  12
        Socket(s):           1
        Stepping:            2
        CPU max MHz:         5000,0000
        CPU min MHz:         800,0000
        BogoMIPS:            7219.20
        Flags:               fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush dts acpi mmx fxsr sse sse2 ss ht tm pbe syscall nx pdpe1gb rdtscp lm constant_tsc art arch_perfmon pebs bt
                             s rep_good nopl xtopology nonstop_tsc cpuid aperfmperf tsc_known_freq pni pclmulqdq dtes64 monitor ds_cpl vmx smx est tm2 ssse3 sdbg fma cx16 xtpr pdcm sse4_1 sse4_2 x2apic movbe popcnt 
                             tsc_deadline_timer aes xsave avx f16c rdrand lahf_lm abm 3dnowprefetch cpuid_fault ssbd ibrs ibpb stibp ibrs_enhanced tpr_shadow flexpriority ept vpid ept_ad fsgsbase tsc_adjust bmi1 avx
                             2 smep bmi2 erms invpcid rdseed adx smap clflushopt clwb intel_pt sha_ni xsaveopt xsavec xgetbv1 xsaves split_lock_detect avx_vnni dtherm ida arat pln pts hwp hwp_notify hwp_act_window h
                             wp_epp hwp_pkg_req hfi vnmi umip pku ospke waitpkg gfni vaes vpclmulqdq tme rdpid movdiri movdir64b fsrm md_clear serialize pconfig arch_lbr ibt flush_l1d arch_capabilities
    Virtualization features: 
      Virtualization:        VT-x
    Caches (sum of all):     
      L1d:                   512 KiB (12 instances)
      L1i:                   512 KiB (12 instances)
      L2:                    12 MiB (9 instances)
      L3:                    25 MiB (1 instance)
    NUMA:                    
      NUMA node(s):          1
      NUMA node0 CPU(s):     0-19
    Vulnerabilities:         
      Gather data sampling:  Not affected
      Itlb multihit:         Not affected
      L1tf:                  Not affected
      Mds:                   Not affected
      Meltdown:              Not affected
      Mmio stale data:       Not affected
      Retbleed:              Not affected
      Spec rstack overflow:  Not affected
      Spec store bypass:     Mitigation; Speculative Store Bypass disabled via prctl
      Spectre v1:            Mitigation; usercopy/swapgs barriers and __user pointer sanitization
      Spectre v2:            Mitigation; Enhanced / Automatic IBRS; IBPB conditional; RSB filling; PBRSB-eIBRS SW sequence; BHI BHI_DIS_S
      Srbds:                 Not affected
      Tsx async abort:       Not affected
      
    $ lshw -short -C memory
    H/W path           Device     Class          Description
    ========================================================
    /0/0                          memory         64KiB BIOS
    /0/3b                         memory         40GiB System Memory
    /0/3b/0                       memory         8GiB DIMM Synchronous 4800 MHz (0,2 ns)
    /0/3b/1                       memory         16GiB DIMM Synchronous 4800 MHz (0,2 ns)
    /0/3b/2                       memory         16GiB DIMM Synchronous 4800 MHz (0,2 ns)
    /0/3b/3                       memory         [empty]
    /0/4b                         memory         384KiB L1 cache
    /0/4c                         memory         256KiB L1 cache
    /0/4d                         memory         10MiB L2 cache
    /0/4e                         memory         25MiB L3 cache
    /0/4f                         memory         128KiB L1 cache
    /0/50                         memory         256KiB L1 cache
    /0/51                         memory         2MiB L2 cache
    /0/52                         memory         25MiB L3 cache
    /0/100/14.2                   memory         RAM memory
```

***UPD***
Есть нюанс, что за время расчета берется время, когда клиент отправит данные с решением, но тут прибавляется ко времени расчета еще и время задержки по сети.
Подразумевается, что задержка будет минимальной и клиент будет коммуницировать с сервером, который максимально близко к нему расположен.
