# Logging Example

## Sample Log Output

When you run the application, `logs.txt` will contain entries like this:

```
[2025-11-04 15:30:01.123] [INFO] Logger initialized successfully
[2025-11-04 15:30:01.124] [INFO] Starting Acts Service...
[2025-11-04 15:30:01.125] [INFO] [METHOD_INIT] ConnectMongoDB
[2025-11-04 15:30:01.126] [INFO] [MONGO_TRANSACTION] CONNECT: Attempting to connect to MongoDB
[2025-11-04 15:30:01.234] [INFO] [MONGO_TRANSACTION] PING: Verifying MongoDB connection
[2025-11-04 15:30:01.245] [INFO] Successfully connected to MongoDB at mongodb://localhost:27017
[2025-11-04 15:30:01.246] [INFO] [METHOD_SUCCESS] ConnectMongoDB
[2025-11-04 15:30:01.247] [INFO] Server starting on 0.0.0.0:8080

--- User creates a new act via API ---

[2025-11-04 15:32:15.100] [INFO] [METHOD_INIT] ActHandler.CreateAct
[2025-11-04 15:32:15.101] [INFO] Received request to create act from IP: 192.168.1.10
[2025-11-04 15:32:15.102] [INFO] [METHOD_INIT] ActService.CreateAct
[2025-11-04 15:32:15.103] [INFO] [METHOD_INIT] ActRepository.Create
[2025-11-04 15:32:15.104] [INFO] [MONGO_TRANSACTION] INSERT: Inserting new act into database
[2025-11-04 15:32:15.234] [INFO] Successfully created act with ID: 673456789abcdef012345678
[2025-11-04 15:32:15.235] [INFO] [METHOD_SUCCESS] ActRepository.Create
[2025-11-04 15:32:15.236] [INFO] Successfully created act with ID: 673456789abcdef012345678
[2025-11-04 15:32:15.237] [INFO] [METHOD_SUCCESS] ActService.CreateAct
[2025-11-04 15:32:15.238] [INFO] Successfully created act via API, ID: 673456789abcdef012345678
[2025-11-04 15:32:15.239] [INFO] [METHOD_SUCCESS] ActHandler.CreateAct

--- User generates Excel document ---

[2025-11-04 15:33:20.100] [INFO] [METHOD_INIT] ActHandler.GenerateAct
[2025-11-04 15:33:20.101] [INFO] Received request to generate act with ID: 673456789abcdef012345678 from IP: 192.168.1.10
[2025-11-04 15:33:20.102] [INFO] [METHOD_INIT] ActService.GenerateAct
[2025-11-04 15:33:20.103] [INFO] Generating act for ID: 673456789abcdef012345678
[2025-11-04 15:33:20.104] [INFO] [METHOD_INIT] ActRepository.FindByID
[2025-11-04 15:33:20.105] [INFO] [MONGO_TRANSACTION] SELECT: Finding act by ID: 673456789abcdef012345678
[2025-11-04 15:33:20.150] [INFO] Successfully found act with ID: 673456789abcdef012345678
[2025-11-04 15:33:20.151] [INFO] [METHOD_SUCCESS] ActRepository.FindByID
[2025-11-04 15:33:20.152] [INFO] Processing and generating act: 673456789abcdef012345678
[2025-11-04 15:33:20.153] [DEBUG] Using 5 positions with current period costs
[2025-11-04 15:33:20.154] [INFO] [METHOD_INIT] ActRepository.Update
[2025-11-04 15:33:20.155] [INFO] [MONGO_TRANSACTION] UPDATE: Updating act with ID: 673456789abcdef012345678
[2025-11-04 15:33:20.200] [INFO] Successfully updated act with ID: 673456789abcdef012345678
[2025-11-04 15:33:20.201] [INFO] [METHOD_SUCCESS] ActRepository.Update
[2025-11-04 15:33:20.202] [INFO] [METHOD_INIT] ExcelService.GenerateAct
[2025-11-04 15:33:20.203] [INFO] [EXCEL_INIT] Starting Excel document creation: generated/act_673456789abcdef012345678_1730736800.xlsx
[2025-11-04 15:33:20.204] [INFO] Opening Excel template: templates/act_template.xlsx
[2025-11-04 15:33:20.350] [INFO] Processing 1 sheets in Excel template
[2025-11-04 15:33:20.351] [DEBUG] Processing sheet: Sheet1
[2025-11-04 15:33:21.100] [INFO] Saving Excel file to: generated/act_673456789abcdef012345678_1730736800.xlsx
[2025-11-04 15:33:21.450] [INFO] [EXCEL_COMPLETE] Excel document created successfully: generated/act_673456789abcdef012345678_1730736800.xlsx
[2025-11-04 15:33:21.451] [INFO] [METHOD_SUCCESS] ExcelService.GenerateAct
[2025-11-04 15:33:21.452] [INFO] [METHOD_INIT] ActRepository.Update
[2025-11-04 15:33:21.453] [INFO] [MONGO_TRANSACTION] UPDATE: Updating act with ID: 673456789abcdef012345678
[2025-11-04 15:33:21.500] [INFO] Successfully updated act with ID: 673456789abcdef012345678
[2025-11-04 15:33:21.501] [INFO] [METHOD_SUCCESS] ActRepository.Update
[2025-11-04 15:33:21.502] [INFO] Successfully generated act with download link: /api/act/download/act_673456789abcdef012345678_1730736800.xlsx
[2025-11-04 15:33:21.503] [INFO] [METHOD_SUCCESS] ActService.GenerateAct
[2025-11-04 15:33:21.504] [INFO] Successfully generated act via API, download link: /api/act/download/act_673456789abcdef012345678_1730736800.xlsx
[2025-11-04 15:33:21.505] [INFO] [METHOD_SUCCESS] ActHandler.GenerateAct

--- User downloads the file ---

[2025-11-04 15:34:00.100] [INFO] [METHOD_INIT] ActHandler.DownloadAct
[2025-11-04 15:34:00.101] [INFO] Received request to download file: act_673456789abcdef012345678_1730736800.xlsx from IP: 192.168.1.10
[2025-11-04 15:34:00.102] [INFO] Sending file to client: act_673456789abcdef012345678_1730736800.xlsx
[2025-11-04 15:34:00.103] [INFO] [METHOD_SUCCESS] ActHandler.DownloadAct

--- Server shutdown ---

[2025-11-04 16:00:00.000] [INFO] Shutting down server...
[2025-11-04 16:00:00.001] [INFO] [METHOD_INIT] MongoDBClient.Disconnect
[2025-11-04 16:00:00.002] [INFO] [MONGO_TRANSACTION] DISCONNECT: Closing MongoDB connection
[2025-11-04 16:00:00.050] [INFO] Disconnected from MongoDB
[2025-11-04 16:00:00.051] [INFO] [METHOD_SUCCESS] MongoDBClient.Disconnect
```

## Error Example

When an error occurs, you'll see detailed error logging:

```
[2025-11-04 15:35:00.100] [INFO] [METHOD_INIT] ActHandler.GenerateAct
[2025-11-04 15:35:00.101] [INFO] Received request to generate act with ID: invalid_id from IP: 192.168.1.10
[2025-11-04 15:35:00.102] [INFO] [METHOD_INIT] ActService.GenerateAct
[2025-11-04 15:35:00.103] [INFO] Generating act for ID: invalid_id
[2025-11-04 15:35:00.104] [INFO] [METHOD_INIT] ActRepository.FindByID
[2025-11-04 15:35:00.105] [ERROR] Invalid ObjectID format: encoding/hex: invalid byte: U+0069 'i'
[2025-11-04 15:35:00.106] [ERROR] [METHOD_ERROR] ActRepository.FindByID: encoding/hex: invalid byte: U+0069 'i'
[2025-11-04 15:35:00.107] [ERROR] [METHOD_ERROR] ActService.GenerateAct: act not found: invalid ID format
[2025-11-04 15:35:00.108] [ERROR] [METHOD_ERROR] ActHandler.GenerateAct: act not found: invalid ID format
```

## Log Analysis

You can easily analyze logs using standard Unix tools:

### View recent logs
```bash
tail -f logs.txt
```

### Count error occurrences
```bash
grep "\[ERROR\]" logs.txt | wc -l
```

### Find all MongoDB operations
```bash
grep "MONGO_TRANSACTION" logs.txt
```

### Find all Excel generation events
```bash
grep "EXCEL_" logs.txt
```

### View logs from specific time
```bash
grep "2025-11-04 15:3" logs.txt
```

### Count method executions
```bash
grep "METHOD_SUCCESS" logs.txt | wc -l
```

### Find failed operations
```bash
grep "METHOD_ERROR" logs.txt
```

## Benefits in Practice

1. **Debugging**: Trace the exact flow of execution through all layers
2. **Performance**: Measure time between init and success events
3. **Monitoring**: Track MongoDB operations and Excel generation
4. **Security**: Audit all API requests with client IPs
5. **Compliance**: Permanent record of all system activities

