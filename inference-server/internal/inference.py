from transformers import AutoModelForSequenceClassification
from transformers import AutoTokenizer
import torch

model_ckpt = "./inference-server/model/distilbert-finetuned-for-classification"
num_labels = 6


tokenizer = AutoTokenizer.from_pretrained(model_ckpt, local_files_only=True)

model = model = AutoModelForSequenceClassification.from_pretrained(model_ckpt, num_labels=num_labels)

def model_output(text):
    input = tokenizer(text, return_tensors="pt", padding=True, truncation=True)
    device = torch.device("cuda" if torch.cuda.is_available() else "cpu")
    model.to(device)
    output = model(**input)
    pred_label = torch.argmax(output.logits, axis=-1)

    return pred_label.item()

