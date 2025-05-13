function api_books_get()
  local method = request.method
  local path = request.url
  local req_body = request.body
  local headers = request.headers

  -- Parsowanie inventoryNumber
  local inventoryNumber = nil
  if req_body then
    local _, _, value = string.find(req_body, "inventoryNumber%s*:%s*(%d+)")
    if value then
      inventoryNumber = tonumber(value)
    end
  end

  -- logs part:
  local log = {}
  table.insert(log, "$ start simulation")
  table.insert(log, "$ received inventoryNumber = " .. (inventoryNumber or "nil"))

  -- tworzenie ciała odpowiedzi jako JSON string
  local response_body = ""
  if not inventoryNumber then
    response_body = [[{ "message": "Książka o takim numerze nie istnieje" }]]
  elseif inventoryNumber == 200 then
    response_body = [[{ "available_inventory_number": "200" }]]
  else
    response_body = [[{ "message": "Brak poprawnego numeru inwentarzowego w ścieżce" }]]
  end

  -- wynik końcowy
  local result = {
    response = {
      status = "ok",
      body = response_body
    },
    log = log
  }

  return result
end
