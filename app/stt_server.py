from fastapi import FastAPI, Response, Request, File, UploadFile

# import basemodel
app = FastAPI()
from faster_whisper import WhisperModel

model = WhisperModel(
    "aslon1213/whisper-small-uz-with-uzbekvoice-ct2",
     device="cuda",
)
import io


@app.post("/")
async def synthesize(request: Request):
    body = await request.body()
    body_io = io.BytesIO(body)
    print("Working with file of length: ", len(body))
    statements, info = model.transcribe(body_io, language="uz")
    full_statement = ""
    for statement in statements:
        print("Statement: ", statement)
        full_statement += statement.text + " "
    return {"text": full_statement, "info": info}
